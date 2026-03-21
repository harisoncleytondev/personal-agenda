import { useState } from "react";

interface AuthPageProps {
  onLogin: () => void;
}

export default function AuthPage({ onLogin }: AuthPageProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!email || !password || (!isLogin && !name)) {
      setError(true);
      setTimeout(() => setError(false), 400);
      return;
    }

    onLogin();
  };

  return (
    <div className="page active auth-page">
      <div className="auth-header">
        <div className="greeting">
          {isLogin ? "Bem-vindo" : "Criar conta"}
          <br />
          <em>{isLogin ? "de volta." : "agora."}</em>
        </div>
      </div>

      <form className="auth-form form-zone" onSubmit={handleSubmit}>
        {!isLogin && (
          <div className="bs-field">
            <label className="bs-label">Nome</label>
            <input
              type="text"
              className={`bs-input ${error && !name ? "err" : ""}`}
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Seu nome"
            />
          </div>
        )}

        <div className="bs-field">
          <label className="bs-label">E-mail</label>
          <input
            type="email"
            className={`bs-input ${error && !email ? "err" : ""}`}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="seu@email.com"
          />
        </div>

        <div className="bs-field">
          <label className="bs-label">Senha</label>
          <input
            type="password"
            className={`bs-input ${error && !password ? "err" : ""}`}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
          />
        </div>

        <button type="submit" className="btn-primary-custom auth-btn">
          {isLogin ? "Entrar" : "Cadastrar"}
        </button>

        <button
          type="button"
          className="btn-secondary-custom auth-toggle"
          onClick={() => {
            setIsLogin(!isLogin);
            setError(false);
          }}
        >
          {isLogin ? "Não tem uma conta? Cadastre-se" : "Já tem conta? Entre"}
        </button>
      </form>
    </div>
  );
}
