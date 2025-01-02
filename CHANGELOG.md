## Currently in dev

### Notes for town administrators

Because of the new chunking size, your reverse proxy may block requests to the endpoints Zap requires. Because of this we've modified our official Nginx configuration as well. You basically just need to add the following in your chat config to make sure the proxy actually handles the thing:

```sh
# Zap upload endpoint (just to make sure nothing happens)
location /auth/liveshare/upload {
  proxy_http_version 1.1;
  client_max_body_size 1100k;
  proxy_pass http://localhost:4001/auth/liveshare/upload;
}
```

That's about it for the breaking changes tough.

### Changes

- Incremented protocol version to v7 (due to the Zap changes)
- Made the registration a little bit more user friendly
  - The display name and username input now have a max length associated with them
  - The display name and username errors for requirements now include the requirements
  - The display name and username creation have been separated to avoid confusion between the two
  - Display name creation now has a better description of what it actually is
  - Username creation now has a better description of what it actually is
- Fixed the email not changing when pressing the resend email button and with a changed email
- Allow a new Zap chunking size for faster performance (512 KB -> 1 MB)
- Made Zap a little faster by increasing the chunks loaded ahead (now from 10 MB max -> 20 MB max)
- Added automatic layering to Tabletop to make playing card games with card stacking easier
- Added new events for Warp to make port sharing a possibility (in Spaces)
- Decentralized connections to Spaces are now possible

## 0.5.1

- Fixed a bug where stored actions would be completely broken (friend requests and stuff)

## 0.5.0

- Incremented protocol version to v6
- Changed the way clients connect to the websocket gateway to get the data through a packet instead of protocols
- Removed LiveKit support from Spaces (read main/README.md for more info)
- Added handlers for actions from Spaces messaging
- Removed some old and unused code

## 0.4.0

### Notes for town administrators

There is a new admin panel in the client app now. Some settings that you could set through the environment file before are now available through the Town page in the settings. Things such as whether decentralization with unsafe locations is allowed or the max file upload size or the maximum amount of storage usage per account can now all be managed through there. I hope this job makes it a little easier to maintain a town. More additions are coming the future, but for now those are all I've got. If you want any further additions be sure to open an issue about it.

### Changes

- Create an invite matching the SYSTEM_UUID if there is no account
- The first account created now instantly gets admin privileges
- Removed unused routes (/app/\* and /node/manage/regen + remove)
- Some code cleanup moving all entities into the database module
- Account management for admins (rank changing, searching, deleting)
- Moved file max size and total storage settings to the admin panel
- New settings for decentralization in the admin panel
- Settings bridge using node tokens as authorization between the nodes and the backend
- Fixed a bug where decentralization wouldn't work because of token synchronization issues
