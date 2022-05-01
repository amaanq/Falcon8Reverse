use rusb::{Context, Device, DeviceHandle, Result, UsbContext};
use std::time::Duration;

const VENDOR_ID: u16 = 0x195D;
const PRODUCT_ID: u16 = 0x6009;

#[derive(Debug)]
struct Endpoint {
    config: u8,
    iface: u8,
    setting: u8,
    address: u8,
}

fn main() -> Result<()> {
    let mut context = Context::new()?;
    let (mut device, mut handle) =
        open_device(&mut context, VENDOR_ID, PRODUCT_ID).expect("Did not find USB device");

    print_device_info(&mut handle)?;

    let endpoints = find_readable_endpoints(&mut device)?;

    for endpoint in &endpoints {
        println!("{:?}", endpoint);
    }

    let endpoint = endpoints
        .first()
        .expect("No Configurable endpoint found on device");

    let has_kernel_driver = match handle.kernel_driver_active(endpoint.iface) {
        Ok(true) => {
            handle.detach_kernel_driver(endpoint.iface)?;
            true
        }
        _ => false,
    };
    println!(
        "has kernel driver? {} {}",
        has_kernel_driver, endpoint.iface
    );

    println!("configuring endpoint");
    configure_endpoint(&mut handle, &endpoint)?;

    // println!("setting idle");
    // set_idle(&mut handle).ok();
    // println!("setting report");
    // set_report(&mut handle)?;
    // println!("reading interrupt");
    // let data = read_interrupt(&mut handle, endpoint.address)?;
    // println!("{:02X?}", &data);

    handle.release_interface(endpoint.iface)?;
    if has_kernel_driver {
        handle.attach_kernel_driver(endpoint.iface)?;
    }

    Ok(())
}

fn open_device<T: UsbContext>(
    context: &mut T,
    vid: u16,
    pid: u16,
) -> Option<(Device<T>, DeviceHandle<T>)> {
    let devices = match context.devices() {
        Ok(d) => d,
        Err(_) => return None,
    };

    for device in devices.iter() {
        let device_desc = match device.device_descriptor() {
            Ok(d) => d,
            Err(_) => continue,
        };

        if device_desc.vendor_id() == vid && device_desc.product_id() == pid {
            match device.open() {
                Ok(handle) => return Some((device, handle)),
                Err(_) => continue,
            }
        }
    }
    None
}

fn print_device_info<T: UsbContext>(handle: &mut DeviceHandle<T>) -> Result<()> {
    let device_desc = handle.device().device_descriptor()?;
    let timeout = Duration::from_secs(1);
    let languages = handle.read_languages(timeout)?;

    println!("Active configuration: {}", handle.active_configuration()?);

    if !languages.is_empty() {
        let language = languages[0];
        println!("Language: {:?}", language);

        println!(
            "Manufacturer: {}",
            handle
                .read_manufacturer_string(language, &device_desc, timeout)
                .unwrap_or("Not Found".to_string())
        );
        println!(
            "Product: {}",
            handle
                .read_product_string(language, &device_desc, timeout)
                .unwrap_or("Not Found".to_string())
        );
        println!(
            "Serial Number: {}",
            handle
                .read_serial_number_string(language, &device_desc, timeout)
                .unwrap_or("Not Found".to_string())
        );
    }
    Ok(())
}

fn find_readable_endpoints<T: UsbContext>(device: &mut Device<T>) -> Result<Vec<Endpoint>> {
    let device_desc = device.device_descriptor()?;
    let mut endpoints = vec![];
    for n in 0..device_desc.num_configurations() {
        let config_desc = match device.config_descriptor(n) {
            Ok(c) => c,
            Err(_) => continue,
        };
        // println!("{:#?}", config_desc);
        for interface in config_desc.interfaces() {
            for interface_desc in interface.descriptors() {
                // println!("{:#?}", interface_desc);
                for endpoint_desc in interface_desc.endpoint_descriptors() {
                    // println!("{:#?}", endpoint_desc);
                    endpoints.push(Endpoint {
                        config: config_desc.number(),
                        iface: interface_desc.interface_number(),
                        setting: interface_desc.setting_number(),
                        address: endpoint_desc.address(),
                    });
                }
            }
        }
    }

    Ok(endpoints)
}

fn configure_endpoint<T: UsbContext>(
    handle: &mut DeviceHandle<T>,
    endpoint: &Endpoint,
) -> Result<()> {
    println!("setting configuration");
    handle.set_active_configuration(endpoint.config)?;
    println!("claiming interface");
    handle.claim_interface(endpoint.iface)?;
    println!("setting alternate setting");
    handle.set_alternate_setting(endpoint.iface, endpoint.setting)
}

fn set_idle<T: UsbContext>(handle: &mut DeviceHandle<T>) -> Result<usize> {
    let timeout = Duration::from_secs(1);
    const REQUEST_TYPE: u8 = 0x21;
    const REQUEST: u8 = 0x0A;
    const VALUE: u16 = 0x0000;
    const INDEX: u16 = 0x0000;
    handle.write_control(REQUEST_TYPE, REQUEST, VALUE, INDEX, &[], timeout)
}

fn set_report<T: UsbContext>(handle: &mut DeviceHandle<T>) -> Result<usize> {
    let timeout = Duration::from_secs(1);

    const REQUEST_TYPE: u8 = 0x21;
    const REQUEST: u8 = 0x09;
    const VALUE: u16 = 0x0307;
    const INDEX: u16 = 0x0002;

    const DATA: [u8; 264] = REPORT_DATA();

    handle.write_control(REQUEST_TYPE, REQUEST, VALUE, INDEX, &DATA, timeout)
}

const fn REPORT_DATA() -> [u8; 264] {
    const HEADER: [u8; 3] = [0x07, 0x82, 0x01];
    let mut data = [0; 264];
    data[0] = HEADER[0];
    data[1] = HEADER[1];
    data[2] = HEADER[2];
    data
}

fn read_interrupt<T: UsbContext>(handle: &mut DeviceHandle<T>, address: u8) -> Result<Vec<u8>> {
    let timeout = Duration::from_secs(1);
    let mut buf = [0u8; 64];

    handle
        .read_interrupt(address, &mut buf, timeout)
        .map(|_| buf.to_vec())
}
