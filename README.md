# frezh

frezh is a tool used to import your [Hello Fresh](https://www.hellofresh.com) recipe cards into a system.

## Running

You can use the `docker-ompose.yml` file provided here to start and run the application. This handles the installation of binaries needed to, for example, support OCR.

You must provide, local to the directory in which you are executing the container, a `.env` file that provides the following:

* `CHATGPT_KEY`: an API key that is funded with requests to ChatGPT
  * This is necessary because the import of images and the subsequent OCR processing is imprecise and ChatGPT provides a convenient means of cleaning up the data before tasking the user with final review and cleanup.
