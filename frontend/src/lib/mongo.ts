import { MongoClient, ServerApiVersion } from 'mongodb'


const uri = import.meta.env.MONGO_URI;
if (!uri) throw new Error("Missing MONGO_URI");

// const uri = process.env.MONGO_URI!;

const client = new MongoClient(uri, {
	serverApi: {
		version: ServerApiVersion.v1,
		strict: true,
		deprecationErrors: true,
	}
})

export const db = client.db("auth")
