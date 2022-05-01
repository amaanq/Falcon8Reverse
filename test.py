import usb.util
import usb.core

# can be changed
BUS_ID = 0x07
DEVICE_ID = 0x06

VENDOR_ID = 0x195D
PRODUCT_ID = 0x6009


# print(usb.core.show_devices())
# get usb 
dev = usb.core.find(idVendor=VENDOR_ID, idProduct=PRODUCT_ID)
print(dev)
print(dev.bDescriptorType)