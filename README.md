# Krell Meter

Growing up in the 70s and 80s, I'm a sucker for analog meters.  There's just something evocative about seeing measurements represented on a physical moving object.

When I realized that you could [drive simple analog voltmeters](https://twitter.com/lcamtuf/status/1114689211385278464) with the PWM output of a microcontroller, I knew immediately what I wanted to do.

## Hardware
- 0-3VDC analog voltmeters.  I used [85C1 meters from eBay](https://www.ebay.com/sch/i.html?_from=R40&_trksid=m570.l1313&_nkw=85C1+3V&_sacat=0)
  (Make sure to get voltmeters though, not ammeters.)
- Arduino Pro Micro 3.3v [SparkFun](https://www.sparkfun.com/products/12587)

Here's the schematic for wiring the two meters to the Pro Micro:
![schem](https://github.com/phreakmonkey/KrellMeter/blob/master/misc/KrellMeter-schem.png)

![meters-rear](https://github.com/phreakmonkey/KrellMeter/blob/master/misc/meters-back.png)

I took an image of the 85C1 faceplate I had and whipped up a quick-and-dirty scale in GIMP to replace it with.  I then printed it on my laser printer and just gued it to the back of the existing faceplate with an Elmer's Glue stick.  The result is pretty convincing:

![meters-front](https://github.com/phreakmonkey/KrellMeter/blob/master/misc/meters-front.png)

faceplate/ contains the faceplate images I used.  Be creative!


## Software

#### arduino/main.cpp

This is the source for the Pro Micro firmware.  It's very simple.  It listens
on the serial port for a string starting with A or B followed by a number and sets the PWM output of pins 5 and 6 respectively to the value of the number.

E.g. "A0" will set the PWM output on pin 5 to 0, whereas "A255" will set the PWM output on pin 5 to 255.

#### client/krellmeter.go

This is the client I run in both Windows 10 and Ubuntu Linux.  It looks for "krellmeter.ini" in the current directory unless you specify a path to an INI file on the command line.  The options are fairly self-explanitory.  The "max" parameter for each meter lets you tweak what PWM output makes the meter read "100%".  

When the client first starts it sweeps from 0-100, stopping momentarily at 25%, 50%, 75%, and 100%.  You can use this to gauge whether you need to increase/decrease "max" for either meter. 

(Alternatively, connect to the Pro Micro from a terminal program and just send A200<enter> and see what it reads.  Adjust the number up or down to find the top for that particular meter.)

One huge caveat:  I'm getting the GPU utlization from the Nvidia "nvidia-smi" client on both platforms.  If you don't have an Nvidia GPU (or the Nvidia client library installed) then this won't work for you. 

TODO: Make the client more configurable for differing configurations.

#### binaries/

Contains pre-compiled binaries for the Linux AMD64 platform and Windows 64-bit platforms.  If you use these be sure to grab the krellmeter.ini file from client/ and taylor it for your configuration.  (E.g. port=COM3 or similar for Windows.)

I used [NSSM](https://nssm.cc/download) to make krellmeter.exe into a service on Windows 10, but any method of running it in the background should work since it doesn't require any special permissions.

Have fun!
