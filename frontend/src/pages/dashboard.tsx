import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { useUser } from "~/query/user";

export default function Dashboard() {
    const router = useRouter();
    const query = useUser();

    if (query.isLoading) {
        return <Loading />
    }

    if (query.isError) {
        router.push("/");
        return <></>;
    }

    if (!query.data) {
        router.push("/");
        return <></>;
    }


    return (
        <>
            <Head>
                <title>Dashboard | Ubiquitodo</title>
            </Head>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <h1 className="text-6xl font-bold mb-8">Welcome to Ubiquitodo!</h1>
                <p className="text-2xl mb-8">
                    You have successfully logged in. {query.data.username}
                </p>
                <Link
                    href="/"
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Back to Login
                </Link>
            </main>
        </>
    );
}


function Loading() {
    return (
        <>
            <main className="flex min-h-screen flex-col items-center justify-center">
                <h1 className="text-6xl font-bold mb-8">Welcome to Ubiquitodo!</h1>
                <p className="text-2xl mb-8">
                    Loading...
                </p>
            </main>
        </>
    );
}
