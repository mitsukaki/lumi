import { useLoaderData } from "react-router-dom";
import { Link } from "react-router-dom";
import * as React from "react";

const api = "http://127.0.0.1:8084/api";
const photoCDN = "https://cdn.lumi.mitsukaki.com/";

export default function AlbumPage() {
    const result = useLoaderData();
    const album = result.album;
    const user = result.user;

    return (
        <div class="h-full">
            <Backdrop
                album_id={album._id}
                photo_id={album.cover_photo.photo_id}
            />
            <div class="container mx-auto pb-8 pt-16 text-white">
                {/* <p class="mx-auto">{JSON.stringify(user)}</p> */}
                <p class="text-center text-2xl text-current">{album.title}</p>
                <p class="text-center text-base text-current">
                    <DateText time={album.date}></DateText>
                </p>
                <p class="text-center text-base text-current">{"@" + user.username}</p>
            </div>
            <PhotoGrid album={album} />
            <div class="text-white">
                <p class="text-center text-base text-current pt-12 pb-8">
                    See more albums from <Link class="text-zinc-500" to={"/u/" + user.username}>{"@" + user.username}</Link>
                </p>
            </div>
        </div>
    );
}

export async function AlbumPageLoader({ params }) {
    const albumId = params.album_id;
    const response = await fetch(`${api}/album/${albumId}`);
    const album = await response.json();

    const userId = album.author_user_id;
    const userResponse = await fetch(`${api}/user/${userId}`);
    const user = await userResponse.json();

    return { 
        album: album,
        user: user
    };
}

class Backdrop extends React.Component {
    getAlbumPhoto(album_id, photo_id) {
        return photoCDN + "albums/" + album_id + "/" + photo_id;
    }

    render() {
        return (
            <div class="backdrop h-full absolute" style={{
                backgroundImage: `url(${this.getAlbumPhoto(
                    this.props.album_id,
                    this.props.photo_id,
                )})`
            }}>
                <div class="backdrop-blur h-full">
                    {this.props.children}
                </div>
            </div>
        );
    }
}

class BufferedPhoto extends React.Component {
    getAlbumPhoto(album_id, photo_id) {
        return photoCDN + "albums/" + album_id + "/" + photo_id;
    }

    render() {
        return (
            <img class="aspect-auto rounded-lg mb-4" src={this.getAlbumPhoto(
                this.props.album_id,
                this.props.photo_id,
            )} />
        );
    }
}

class PhotoGrid extends React.Component {
    render() {
        let album_id = this.props.album._id;
        return (
            <div class="container mx-auto">
                {
                    this.props.album.photos.map(photo => (
                        <BufferedPhoto
                            album_id={album_id}
                            photo_id={photo.photo_id}
                        />
                    ))
                }
            </div>
        );
    }
}

class DateText extends React.Component {
    render() {
        const date = new Date(parseInt(this.props.time));
        const daysOfWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
        const monthsOfYear = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

        const dayOfWeek = daysOfWeek[date.getDay()];
        const month = monthsOfYear[date.getMonth()];

        const finalDate =
            dayOfWeek + " " +
            month + " " +
            date.getDate() + ", " +
            date.getFullYear() + " at " +
            (date.getHours() % 12) + ":" +
            date.getMinutes() +
            (date.getHours() > 12 ? " PM" : " AM");
            
        return (
            <span>{finalDate}</span>
        );
    }
}