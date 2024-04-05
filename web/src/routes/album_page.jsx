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
                    See more albums from <Link class="text-zinc-500" to={"/u/" + album.author_user_id }>{"@" + user.username}</Link>
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
                backgroundImage: `url(
                    ${this.getAlbumPhoto(
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
            <img class={`aspect-auto rounded-lg mb-4`} src={this.getAlbumPhoto(
                this.props.album_id,
                this.props.photo_id,
            )} />
        );
    }
}

class PhotoGrid extends React.Component {
    cols(length) {
        if (length <= 4) return length;

        // special/ideal cases
        else if (length == 8) return 4;

        // general "catch-all" break points
        else if (length % 2 == 0) return 2;
        else if (length % 3 == 0) return 3;
        else if (length % 5 == 0) return 5;

        else return 4;
    }

    smallCols(length) {
        let columnCount = this.cols(length);
        return (columnCount >= 2) ? 2 : columnCount;
    }

    render() {
        let album_id = this.props.album._id;

        // sort the photos into arrays using their "row" property
        let rows = [];
        this.props.album.photos.forEach(photo => {
            if (rows[photo.row] === undefined)
                rows[photo.row] = [];
            
            rows[photo.row].push(photo);
        });

        return (
            <div class="container mx-auto"> {
            rows.map((row, index) => (
                <div class={`grid grid-cols-${this.smallCols(row.length)} md:grid-flow-col md:grid-cols-${this.cols(row.length) } gap-4`}> {
                row.map((photo) => (
                    <BufferedPhoto
                        album_id={album_id}
                        photo_id={photo.photo_id}
                        key={photo.photo_id}
                    />
                ))
                }</div>
            ))
            }</div>
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