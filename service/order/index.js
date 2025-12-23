const express = require("express");
const app = express();

app.get("/health", (req, res) => res.send("order ok"));

app.post("/order", (req, res) => {
  res.json({ service: "order", status: "created" });
});

app.listen(3001, () => console.log("order running"));
