## 0.6.1

- Fixed the mail server not using the identity specified in the environment file/config
- Added from and to headers to the email message

## 0.6.0

### Notes for town administrators

Zap will probably stop working after this update if you used the official Docker tutorial for installation. We have a guide on how to fix it over at https://docs.liphium.com/migration. Please follow the guides to also resolve the other breaking changes.

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
- Fixed subscriptions taking long when some servers are offline
- Fixed the profile picture not loading in a decentralized context

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
