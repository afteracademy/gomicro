function seed(dbName, user, password) {
  db = db.getSiblingDB(dbName);
  db.createUser({
    user: user,
    pwd: password,
    roles: [{ role: "readWrite", db: dbName }],
  });
}

seed("blog-db", "blog-db-user", "changeit");
seed("blog-test-db", "blog-test-db-user", "changeit");
