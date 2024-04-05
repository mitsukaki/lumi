import * as React from "react";
import * as ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import "./index.css";

import Root from "./routes/root";
import AlbumPage from "./routes/album_page";
import ProfilePage from "./routes/profile_page";
import AuthPage from "./routes/auth_page";
import ErrorPage from "./error_page";
import { AlbumPageLoader } from "./routes/album_page";
import { ProfilePageLoader } from "./routes/profile_page";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "auth",
        element: <AuthPage />,
      },
      {
        path: "a/:album_id",
        element: <AlbumPage />,
        loader: AlbumPageLoader,
      },
      {
        path: "u/:user_id",
        element: <ProfilePage />,
        loader: ProfilePageLoader,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);