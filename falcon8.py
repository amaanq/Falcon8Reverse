from typing import List
import usb

VENDOR_ID = 0x195D   # Itron Technology iONE
PRODUCT_ID = 0x6009  # Unknown


class Falcon8:
    def __init__(self):
        self.idVendor = VENDOR_ID
        self.idProduct = PRODUCT_ID
        self.device: usb.Device = usb.core.find(
            idVendor=self.idVendor, idProduct=self.idProduct)
        self.interfaces = [0, 1, 2]

        if self.device is None:
            raise ValueError("You do not have a Falcon 8 plugged in!")

    def detach_kernel_driver(self):
        for interface in self.interfaces:
            if self.device.is_kernel_driver_active(interface) is True:
                print("has kernel driver", interface)
                # tell the kernel to detach
                self.device.detach_kernel_driver(interface)
               # claim the device
            usb.util.claim_interface(self.device, interface)

    def reattach_kernel_driver(self):
        for interface in self.interfaces:
            usb.util.release_interface(self.device, interface)
            # reattach the device to the OS kernel
            self.device.attach_kernel_driver(interface)

    def set_idle(self):
        REQUEST_TYPE = 0x21
        REQUEST = 0x0A
        VALUE = 0x0000
        INDEX = 0x0000

        print("IDLE RESPONSE:", self.device.ctrl_transfer(
            REQUEST_TYPE, REQUEST, VALUE, INDEX, None, 1000))

    def set_report_init(self):
        REQUEST_TYPE = 0x21
        REQUEST = 0x09
        VALUE = 0x0200
        INDEX = 0x0000

        DATA = [0x00]

        print("REPORT RESPONSE", self.device.ctrl_transfer(
            REQUEST_TYPE, REQUEST, VALUE, INDEX, DATA, 1000))

    def set_report_keys(self):
        REQUEST_TYPE = 0x21
        REQUEST = 0x09
        VALUE = 0x0307
        INDEX = 0x0002

        DATA = [0x00]*264
        DATA[0] = 0x07
        DATA[1] = 0x82
        DATA[2] = 0x01

        print("REPORT RESPONSE", self.device.ctrl_transfer(
            REQUEST_TYPE, REQUEST, VALUE, INDEX, DATA, 1000))

    def hid_set_report(self, report):
        """ Implements HID SetReport via USB control transfer """
        REQUEST_TYPE = 0x21
        REQUEST = 0x09
        VALUE = 0x0307
        INDEX = 0x0002
        return self.device.ctrl_transfer(
            REQUEST_TYPE,  # REQUEST_TYPE_CLASS | RECIPIENT_INTERFACE | ENDPOINT_OUT
            REQUEST,       # SET_REPORT
            VALUE,         # "Vendor" Descriptor Type + 0 Descriptor Index
            INDEX,         # USB interface № 0
            report         # the HID payload as a byte array -- e.g. from struct.pack()
        )

    def hid_get_report(self):
        """ Implements HID GetReport via USB control transfer """
        REQUEST_TYPE = 0xA1
        REQUEST = 0x01
        VALUE = 0x0200
        INDEX = 0x0002
        return self.device.ctrl_transfer(
            REQUEST_TYPE,  # REQUEST_TYPE_CLASS | RECIPIENT_INTERFACE | ENDPOINT_IN
            REQUEST,       # GET_REPORT
            VALUE,         # "Vendor" Descriptor Type + 0 Descriptor Index
            INDEX,         # USB interface № 0
            264            # max reply size
        )


if __name__ == "__main__":
    falcon8 = Falcon8()
    falcon8.detach_kernel_driver()
    data = falcon8.hid_get_report()
    print("DATA", data)
    # try:
    #     falcon8.set_idle()
    # except Exception as e:
    #     print(repr(e))

    # try:
    #     falcon8.set_report()
    # except Exception as e:
    #     print(repr(e))

    # print(data[134], data[135])
    # data[134] = 0x04
    # data[135] = 0x01

    # falcon8.hid_set_report(data.tobytes())

    falcon8.reattach_kernel_driver()
