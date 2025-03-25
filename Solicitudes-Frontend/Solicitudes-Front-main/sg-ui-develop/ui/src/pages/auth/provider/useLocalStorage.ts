import { AutenticarToken, Token } from "../types";

const keyName = "user";

export const getUser = (): Token | null => {
  const value = window.localStorage.getItem(keyName);
  if (value && typeof value !== "undefined" && value !== undefined) {
    return JSON.parse(value) as Token;
  }

  return null;
};

export const getAutenticarInfo = (): AutenticarToken | null => {
  const value = window.localStorage.getItem("autenticar");
  if (value && typeof value !== "undefined" && value !== undefined) {
    return JSON.parse(value) as AutenticarToken;
  }

  return null;
};

export const clearLocalStorage = () => {
  sessionStorage.clear();
  localStorage.clear();
};

export const setLocalStorage = (newValue: Token | null) => {
  if (newValue !== undefined) {
    window.localStorage.setItem(keyName, JSON.stringify(newValue));
  }
};

export const setAutenticarStorage = (newValue: Token | null) => {
  if (newValue !== undefined) {
    window.localStorage.setItem("autenticar", JSON.stringify(newValue));
  }
};
