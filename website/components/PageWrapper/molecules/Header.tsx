import React, { useEffect, useState } from "react";
import routes from "common/routes";
import Link from "next/link";
import * as Styles from "../PageWrapper.styles";
import { useRouter } from "next/router";
import { SearchBar } from "./SearchBar";
import { slide as Menu } from "react-burger-menu";
import { checkMobileBreakPoint } from "common/util";

export const Header: React.FC = () => {
  const [isMobile, setIsMobile] = useState(true);
  const router = useRouter();

  const resizeListener = (): void => {
    setIsMobile(checkMobileBreakPoint());
  };

  useEffect(() => {
    // Resize on first render
    resizeListener();

    window.addEventListener("resize", resizeListener);
    return () => window.removeEventListener("resize", resizeListener);
  }, []);

  return isMobile ? (
    <Styles.Mobile.Header>
      <SearchBar />

      <Menu noOverlay>
        <Styles.Title>Filter It!</Styles.Title>
        {Object.values(routes).map((route) => (
          <Styles.Mobile.LinkButton
            key={route.path}
            selected={route.path === router.pathname}
          >
            <Link href={route.path}>{route.displayName}</Link>
          </Styles.Mobile.LinkButton>
        ))}
      </Menu>
    </Styles.Mobile.Header>
  ) : (
    <Styles.HeaderContainer>
      <Styles.Title>Filter It!</Styles.Title>

      <SearchBar />

      <Styles.LinkContainer>
        {Object.values(routes).map((route) => (
          <Styles.LinkButton
            key={route.path}
            selected={route.path === router.pathname}
          >
            <Link href={route.path}>{route.displayName}</Link>
          </Styles.LinkButton>
        ))}
      </Styles.LinkContainer>
    </Styles.HeaderContainer>
  );
};
