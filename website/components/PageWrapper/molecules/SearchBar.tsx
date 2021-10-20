import { Search } from "@material-ui/icons";
import React from "react";
import theme from "styles/colors";
import * as Styles from "../PageWrapper.styles";

export const SearchBar: React.FC = () => {
  return (
    <>
      <Search htmlColor={theme.border} />
      <Styles.SearchInput placeholder="#sepia, #greyscale" />
    </>
  );
};
