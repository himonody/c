// @ts-ignore
import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        success: {
          DEFAULT: "hsl(var(--success))",
          foreground: "hsl(var(--success-foreground))",
        },
        warning: {
          DEFAULT: "hsl(var(--warning))",
          foreground: "hsl(var(--warning-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
      },
      borderRadius: {
        sm: "var(--radius-sm)",    // 12px
        DEFAULT: "var(--radius-md)", // 16px
        md: "var(--radius-md)",      // 16px
        lg: "var(--radius-lg)",      // 24px
        xl: "var(--radius-lg)",      // 24px
        "2xl": "var(--radius-lg)",   // 24px
        "3xl": "var(--radius-lg)",   // 24px
      },
      maxWidth: {
        container: "1200px",
      },
      height: {
        "14": "56px", // 主按钮高度
      },
      boxShadow: {
        card: "0 8px 24px rgba(0, 0, 0, 0.06)",
        hover: "0 12px 32px rgba(0, 0, 0, 0.08)",
      },
    },
  },
  plugins: [],
};
export default config;
