import { useState } from "react";
import { API_URL } from "./App.tsx";

export const CreatePostForm: React.FC = () => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  const handleSubmit = async () => {
    await fetch(`${API_URL}/posts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify({
        title,
        content,
      }),
    });

    setTitle("");
    setContent("");
  }

  return (
    <div className="gopher-form">
      <label>
        <input
          placeholder="Title"
          value={title}
          type="text"
          onChange={(e) => setTitle(e.target.value)}
        />
      </label>
      <label>
        <textarea
          placeholder="What's in your mind?"
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
      </label>
      <button onClick={handleSubmit}>Create Post</button>
    </div>
  )
}