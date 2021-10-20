import { Filter } from "api/types";

export const filterStringMap: Record<Filter, string> = {
  [Filter.GRAYSCALE]: "#greyscale",
  [Filter.INVERT]: "#invert",
  [Filter.SEPIA]: "#sepia",
};
