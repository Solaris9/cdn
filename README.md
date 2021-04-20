# New Document# CDN

This is my own personal CDN written in GO and Svelte using DigitalOceans Spaces.

## Credits

Anthony's repository as I would have been lost with DigitalOcean, Firebase and other things without it: <https://github.com/acollierr17/cdn>

## Building

This uses Make to make it easier to build the frontend and backend into a folder. \
You are not required to use Make if that's what you prefer, simply execute commands in the `build` rule.

### Client

Build the frontend using `yarn build` inside of the `client` folder, it'll output the built files into `public/build`

### Server

Golang by default builds to the current operating system, if you would like to build for Windows or Linux then you need to set the `GOOS` environment variable (check [GOOS examples](#goos-examples)). \
You can find more information about this here: <https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures>

Then build the executable using `go build`, by default it'll output in the current directory, if you want it elsewhere you'll need to use the `-o` option with a directory.

#### GOOS Examples

Unix:

```bash
GOOS=linux
```

Windows PowerShell:

```bash
$env:GOOS = "linux"
```

### Firebase, DigitalOcean and server configuration

To connect to a Firebase App you will need the authentication file `service-account.json` in the current directory. \
For DigitalOcean Spaces you will need to add the necessary info that start with `SPACES` in the [.env](/.env.example) file, which should also be in the current directory.

More explanation of the rest of the environment variables:

`CDN_ENDPOINT` is your site endpoint, such as `https://cdn.mysite.com` \
`AUTHORIZATION` is the main authorization token, this should be kept as anyone will be able to upload and delete files through the site.

## Todo

- [x] Server routes
  - [x] Files
    - [x] Upload
    - [x] Delete
    - [x] Retrieve
  - [x] Folders
    - [x] Create
    - [x] Edit
    - [x] Retrieve
    - [x] Delete
  - [ ] Get user info such as amount of files, total size
  - [ ] Live socket

- [x] Dashboard
  - [x] View and manage files
    - [ ] Upload directly
    - [x] Delete files
    - [ ] View files?
  - [ ] View and manage folders
    - [ ] Create folders
    - [ ] Edit folders, add/remove files
    - [ ] Retrieve folders and their files
    - [ ] Delete

## Planned things

- [ ] Authentication with firebase auth
  - [ ] Multiple users
  - [ ] Admin account
