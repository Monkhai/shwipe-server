# TODO

- syncronize DB with sessions and users
  - sessions should be written to DB and removed when closed
    - add users to session_users table and remove when session/user is removed.
  - users should be taken from DB and not from Google to unify the types
- create getFriendsSessions endpoint.
