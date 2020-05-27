# Example of low level interaction with a BLE UART device that has an RX and TX
# characteristic for receiving and sending data.  This doesn't use any service
# implementation and instead just manipulates the services and characteristics
# on a device.  See the uart_service.py example for a simpler UART service
# example that uses a high level service implementation.
# Author: Tony DiCola
import logging
import time
import uuid

import Adafruit_BluefruitLE


# Enable debug output.
#logging.basicConfig(level=logging.DEBUG)

# Define service and characteristic UUIDs used by the Automation IO service.
AUTOMATION_IO_SERVICE_UUID = uuid.UUID('00001815-0000-1000-8000-00805F9B34FB')
ANALOG_CHAR_UUID = uuid.UUID('00002A58-0000-1000-8000-00805F9B34FB')

# Get the BLE provider for the current platform.
ble = Adafruit_BluefruitLE.get_provider()

# Main function implements the program logic so it can run in a background
# thread.  Most platforms require the main thread to handle GUI events and other
# asyncronous events like BLE actions.  All of the threading logic is taken care
# of automatically though and you just need to provide a main function that uses
# the BLE provider.
def main():
    # Clear any cached data because both bluez and CoreBluetooth have issues with
    # caching data and it going stale.
    ble.clear_cached_data()

    # Get the first available BLE network adapter and make sure it's powered on.
    adapter = ble.get_default_adapter()
    adapter.power_on()
    print('Using adapter: {0}'.format(adapter.name))

    # Disconnect any currently connected Automation IO devices.
    # Good for cleaning up and starting from a fresh state.
    print('Disconnecting any connected Automation IO devices...')
    ble.disconnect_devices([AUTOMATION_IO_SERVICE_UUID])

    # Scan for Automationm IO devices.
    print('Searching for Automation IO device...')
    try:
        adapter.start_scan()
        # Search for the first Automation IO device found (will time out after 60 seconds
        # but you can specify an optional timeout_sec parameter to change it).
        device = ble.find_device(service_uuids=[AUTOMATION_IO_SERVICE_UUID])
        if device is None:
            raise RuntimeError('Failed to find Automation IO device!')
    finally:
        # Make sure scanning is stopped before exiting.
        adapter.stop_scan()

    print('Connecting to device...')
    device.connect()  # Will time out after 60 seconds, specify timeout_sec parameter
                      # to change the timeout.

    # Once connected do everything else in a try/finally to make sure the device
    # is disconnected when done.
    try:
        # Wait for service discovery to complete for at least the specified
        # service and characteristic UUID lists.  Will time out after 60 seconds
        # (specify timeout_sec parameter to override).
        print('Discovering services...')
        device.discover([AUTOMATION_IO_SERVICE_UUID], [ANALOG_CHAR_UUID])

        # Find the Automation IO service and its characteristics.
        automation = device.find_service(AUTOMATION_IO_SERVICE_UUID)
        analog = automation.find_characteristic(ANALOG_CHAR_UUID)

        # Function to receive Analog characteristic changes.  Note that this will
        # be called on a different thread so be careful to make sure state that
        # the function changes is thread safe.  Use queue or other thread-safe
        # primitives to send data to other threads.
        def received(data):
            value = int.from_bytes(data, byteorder='little')
            print('Received: {0}'.format(value))

        # Turn on notification of Analog characteristics using the callback above.
        print('Subscribing to Analog characteristic changes...')
        analog.start_notify(received)

        # Wait forever. Notify handling happens in another thread
        while True:
          time.sleep(1)
    finally:
        # Make sure device is disconnected on exit.
        device.disconnect()


# Initialize the BLE system.  MUST be called before other BLE calls!
ble.initialize()

# Start the mainloop to process BLE events, and run the provided function in
# a background thread.  When the provided main function stops running, returns
# an integer status code, or throws an error the program will exit.
ble.run_mainloop_with(main)
