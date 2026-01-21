function seed(dbName, user, password) {
  db = db.getSiblingDB(dbName);
  db.createUser({
    user: user,
    pwd: password,
    roles: [{ role: "readWrite", db: dbName }],
  });

  db.createCollection("dummy");

  db.dummy.insert({
    value: "This is a dummy document",
  });
}

seed("blog-db", "blog-db-user", "changeit");
seed("blog-test-db", "blog-test-db-user", "changeit");
