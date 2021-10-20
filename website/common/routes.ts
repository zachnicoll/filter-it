interface Route {
  path: string;
  displayName: string;
}

const routes: Record<string, Route> = {
  Feed: { path: "/", displayName: "Feed" },
  Upload: { path: "/upload", displayName: "Upload" },
};

export default routes;
