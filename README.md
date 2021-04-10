# CDN

This is my own personal CDN written in GO and Svelte using DigitalOceans Spaces, though it'll use Firebase's authentication so anyone will be able to use it if they wish to.

## Credits

Anthony's repository as I would have been lost with DigitalOcean, Firebase and other things without it: <https://github.com/acollierr17/cdn>

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

- [ ] Dashboard
  - [ ] View and manage files
    - [ ] Upload directly
    - [ ] Delete files
    - [ ] View files
  - [ ] View and manage folders
    - [ ] Create folders
    - [ ] Edit folders, add/remove files
    - [ ] Retrieve folders and their files
    - [ ] Delete

- [ ] Authentication with firebase auth
  - [ ] Multiple users
  - [ ] Admin account
