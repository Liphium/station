## Currently in dev

- Made the registration a little bit more user friendly
  - The display name and username input now have a max length associated with them
  - The display name and username errors for requirements now include the requirements
  - The display name and username creation have been separated to avoid confusion between the two
  - Display name creation now has a better description of what it actually is
  - Username creation now has a better description of what it actually is
- Fixed the email not changing when pressing the resend email button and with a changed email
- Allows for a new Zap chunking size for faster performance

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