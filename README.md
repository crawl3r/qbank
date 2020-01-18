# qbank
A little question and answer cli thing I made to help with my CPSA and general knowledge. I have finished with it now but hopefully it might help others. I created it from internet findings, already known topics and friends input.

Current question topics:
- acronyms
- ports
- webservers

More can be added, just follow the structure of the json files and make sure it is in the data/questions/ directory. Anything in there should load fine if it follows the same layout.

Building: Usual Golang stuff, can be run in the Go environment or built to a target system of your choice.

Known issues? 
* Windows sometimes doesn't like the input or line endings. I have only used this on OSX, but feel free to PR a Windows fix if anyone cleans it up.
* SOME questions might have incorrect answers, a few were pointed out to my by a friend but I forgot which ones they were. If you see them and change them, PR them up :)
