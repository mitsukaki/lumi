import { useLoaderData } from "react-router-dom";
import { useState, useEffect } from "react";

const api = "http://127.0.0.1:8084/api";

export default function CreatePage() {
    const [currentDate, setCurrentDate] = useState("");

    useEffect(() => {
        const today = new Date();
        const formattedDate = today.toISOString().split("T")[0];
        setCurrentDate(formattedDate);
    }, []);

    return (
        <div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
            <div class="sm:mx-auto sm:w-full sm:max-w-sm">
                <img class="mx-auto h-10 w-auto" src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600" alt="Your Company" />
                <h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">Create Album</h2>
            </div>
            <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                <form className="space-y-6" onSubmit={createAlbum}>
                    <div className="mb-4">
                        <label htmlFor="title" className="block text-lg font-medium text-gray-700">Title</label>
                        <input type="text" id="title" name="title" className="mt-1 p-2 border border-gray-300 rounded-md w-full" />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="description" className="block text-lg font-medium text-gray-700">Description</label>
                        <textarea id="description" name="description" className="mt-1 p-2 border border-gray-300 rounded-md w-full" />
                    </div>
                    <div className="mb-4">
                        <label htmlFor="date" className="block text-lg font-medium text-gray-700">Date</label>
                        <input type="date" id="date" name="date" className="mt-1 p-2 border border-gray-300 rounded-md w-full" defaultValue={currentDate} />
                    </div>
                    <div className="mb-4">
                        <input type="checkbox" id="unlisted" name="unlisted" className="mr-2" />
                        <label htmlFor="unlisted" className="text-lg font-medium text-gray-700">Unlisted</label>
                    </div>
                    <div>
                        <button type="submit" class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Create</button>
                    </div>
                </form>
            </div>
        </div>
    );
}

function createAlbum(e) {
    e.preventDefault();

    // get user id from local storage
    const user_id = localStorage.getItem("user_id");

    // Get form data
    const title = e.target.elements.title.value;
    const description = e.target.elements.description.value;
    const date = e.target.elements.date.value;
    const unlisted = e.target.elements.unlisted.checked;

    // Create album
    fetch(`${api}/user/${user_id}/album`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            description: description,
            date: date,
            title: title,
            private: unlisted
        })
    }).then(response => {
        if (response.ok) {
            // redirect to image upload page
            response.json().then(data => {
                console.log(data);
                window.location.href = "/upload/" + data.album_id;
            });
        } else {
            console.log("Failed to create album");
        }
    }).catch(error => {
        console.log(error);
    });
}
