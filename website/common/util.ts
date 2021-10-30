import { Filter } from "api/types";
import { filterStringMap } from "./enum.maps";

export const checkMobileBreakPoint = (): boolean => {
  return window.innerWidth <= 600;
};

interface FilterRadioButton {
  label: string;
  value: Filter;
}

export const filterRadioButtons: FilterRadioButton[] = [
  {
    label: filterStringMap[Filter.GRAYSCALE],
    value: Filter.GRAYSCALE,
  },
  {
    label: filterStringMap[Filter.SEPIA],
    value: Filter.SEPIA,
  },
  {
    label: filterStringMap[Filter.INVERT],
    value: Filter.INVERT,
  },
];
