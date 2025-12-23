const express = require("express");
const app = express();

app.get("/health", (req, res) => {
  res.status(200).send("Auth service healthy 🚀");
});

app.get("/login", (req, res) => {
  res.json({ message: "Login API working" });
});

app.listen(3000, () => {
  console.log("Auth service running on port 3000");
});
