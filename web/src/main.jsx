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

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
  },
  {
    path: "auth",
    element: <AuthPage />
  },
  {
    path: "album/:user_id",
    element: <AlbumPage />,
  },
  {
    path: "profile/:user_id",
    element: <ProfilePage />,
  }
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);