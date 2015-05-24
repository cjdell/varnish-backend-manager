
if (req.http.host == "a.example.com") {
  set req.backend_hint = default;
}
	
if (req.http.host == "b.example.com") {
  set req.backend_hint = default;
}
	
if (req.http.host == "c.example.com") {
  set req.backend_hint = default;
}
	