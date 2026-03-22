import { useState } from "react";
import { GET_BASE_URL } from "../utils/constants";

interface AuthPageProps {
  onLogin: () => void;
}

export default function AuthPage({ onLogin }: AuthPageProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordCo, setPasswordCo] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email || !password || (!isLogin && (!name || !passwordCo))) {
      setError(true);
      setTimeout(() => setError(false), 400);
      return;
    }

    setLoading(true);
    try {
      if (isLogin) {
        const res = await fetch(`${GET_BASE_URL}/auth/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password }),
          credentials: "include",
        });
        if (res.ok) {
          onLogin();
        } else {
          setError(true);
          setTimeout(() => setError(false), 400);
        }
      } else {
        const res = await fetch(`${GET_BASE_URL}/auth/register`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name,
            email,
            password,
            password_confirm: passwordCo,
          }),
        });
        if (res.ok) {
          setIsLogin(true);
          setPassword("");
          setPasswordCo("");
        } else {
          setError(true);
          setTimeout(() => setError(false), 400);
        }
      }
    } catch (err) {
      setError(true);
      setTimeout(() => setError(false), 400);
    } finally {
      setLoading(false);
    }
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

        {!isLogin && (
          <div className="bs-field">
            <label className="bs-label">Confirmar Senha</label>
            <input
              type="password"
              className={`bs-input ${error && !passwordCo ? "err" : ""}`}
              value={passwordCo}
              onChange={(e) => setPasswordCo(e.target.value)}
              placeholder="••••••••"
            />
          </div>
        )}

        <button
          type="submit"
          className="btn-primary-custom auth-btn"
          disabled={loading}
        >
          {loading ? "Carregando..." : isLogin ? "Entrar" : "Cadastrar"}
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
