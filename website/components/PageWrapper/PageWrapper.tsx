import React from "react";
import * as Styles from "./PageWrapper.styles";
import { Header } from "./molecules/Header";

export const PageWrapper: React.FC = ({ children }) => {
  return (
    <Styles.Container>
      <Header />
      {children}
    </Styles.Container>
  );
};
