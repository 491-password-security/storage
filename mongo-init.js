db.createUser(
        {
            user: "event",
            pwd: "event",
            roles: [
                {
                    role: "readWrite",
                    db: "event"
                }
            ]
        }
);