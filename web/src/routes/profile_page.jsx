import { useLoaderData } from "react-router-dom";

const api = "http://127.0.0.1:8084/api";

export default function ProfilePage() {
    const result = useLoaderData();
    const user = result.user;

    return (
        <div>
            <h1>Profile</h1>
            <p>Profile page content goes here.</p>
        </div>
    );
}

export async function ProfilePageLoader({ params }) {
    const userId = params.user_id;
    const response = await fetch(`${api}/user/${userId}`);
    const user = await response.json();
    return { user };
}
