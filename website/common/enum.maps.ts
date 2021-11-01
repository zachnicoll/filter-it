import { Filter } from "api/types";

export const filterStringMap: Record<Filter, string> = {
  [Filter.GRAYSCALE]: "#grayscale",
  [Filter.INVERT]: "#invert",
  [Filter.SEPIA]: "#sepia",
};
