import API from "api";
import { Progress } from "api/types";
import { toastError, toastSuccess } from "common/toast";
import {
  Dispatch,
  createContext,
  useEffect,
  useReducer,
  useRef,
  useContext,
} from "react";
import progressReducer, {
  ProgressAction,
  ProgressContextState,
} from "./reducer";
import { useSearch } from "../Search/SearchContext";

interface ProgressContextType {
  dispatchProgress: Dispatch<ProgressAction>;
  progressState: ProgressContextState;
}

const defaultState: ProgressContextState = {};

const ProgressContext = createContext<ProgressContextType | undefined>(
  undefined
);

const ProgressProvider: React.FC = ({ children }) => {
  const [progressState, dispatchProgress] = useReducer(
    progressReducer,
    defaultState
  );
  const { dispatchSearch } = useSearch();

  const intervalRef = useRef<any>(undefined);

  const checkImageProgress = async (): Promise<void> => {
    if (progressState.id) {
      let shouldClearState = false;

      const progress = await API.progress.post(progressState.id);

      if (progress.progress === Progress.DONE) {
        toastSuccess(`Image processed successfully!`);
        shouldClearState = true;
        dispatchSearch({
          type: "SEARCH",
          payload: null,
        });
      } else if (progress.progress === Progress.FAILED) {
        toastError(
          "Failed to process image, please try again",
          new Error("Image Processing Failed")
        );
        shouldClearState = true;
      }

      if (shouldClearState) {
        dispatchProgress({ type: "CLEAR" });
        intervalRef.current && clearInterval(intervalRef.current);
      }
    }
  };

  useEffect(() => {
    if (progressState.id && progressState.status == Progress.PROCESSING) {
      intervalRef.current = setInterval(() => {
        checkImageProgress();
      }, 5000);
    }

    return () => intervalRef.current && clearInterval(intervalRef.current);
  }, [progressState]);

  return (
    <ProgressContext.Provider value={{ progressState, dispatchProgress }}>
      {children}
    </ProgressContext.Provider>
  );
};

const useProgress = (): ProgressContextType => {
  const context = useContext(ProgressContext);
  if (context === undefined) {
    throw new Error("useProgress must be used within a ProgressProvider");
  }
  return context;
};

export { ProgressProvider, useProgress };
