import { betterAuth } from "better-auth";
import { mongodbAdapter } from "better-auth/adapters/mongodb";
import { db } from "./mongo.ts";
import { magicLink } from "better-auth/plugins";

export const auth = betterAuth({
	database: mongodbAdapter(db),
	plugins: [
        magicLink({
            sendMagicLink: async ({ email, url }, request) => {
                // send email to user
            },
            expiresIn: 300,
            disableSignUp: false,
        })
    ]
});
