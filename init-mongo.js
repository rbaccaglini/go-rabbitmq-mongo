// Criar usuário específico para a aplicação
db.createUser({
  user: "appuser",
  pwd: "apppassword",
  roles: [
    {
      role: "readWrite",
      db: "userdb"
    }
  ]
});

// Criar collection inicial
db.createCollection("processed_users");