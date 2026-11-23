import { StrictMode } from 'react'
import './index.css'
import App from './App.tsx'
import { ConfirmationPage } from "./ConfirmationPage.tsx";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { CreatePostForm } from "./CreatePostForm.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/confirm/:token",
    element: <ConfirmationPage />,
  },
  {
    path: "/create-post",
    element: <CreatePostForm />,
  }
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
