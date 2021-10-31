import React from "react";
import * as Styles from "../PageWrapper.styles";
import { Filter } from "api/types";
import { useSearch } from "context";
import theme from "styles/colors";

const options: { value: Filter | null; label: string }[] = [
  { value: Filter.GRAYSCALE, label: "#grayscale" },
  { value: Filter.INVERT, label: "#invert" },
  { value: Filter.SEPIA, label: "#sepia" },
  { value: null, label: "#all" },
];

export const SearchBar: React.FC = () => {
  const { searchState, dispatchSearch } = useSearch();

  return (
    <Styles.Picker
      placeholder={"#filter"}
      options={options}
      styles={{
        control: (styles) => ({
          ...styles,
          flex: 1,
          marginRight: "0.5em",
          borderColor: theme.border,
        }),
      }}
      value={searchState.search}
      onChange={(selected: unknown) => {
        dispatchSearch({
          type: "SEARCH",
          payload: selected,
        });
      }}
    />
  );
};
