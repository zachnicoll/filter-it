import API from "api";
import { Progress } from "api/types";
import { toastError, toastShow, toastSuccess, toastUpdate } from "common/toast";
import {
  Dispatch,
  createContext,
  useEffect,
  useReducer,
  useRef,
  useContext,
} from "react";
import { toast } from "react-toastify";
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
        toastShow('Reload the window to see your filtered image!')

        shouldClearState = true;
      } else if (progress.progress === Progress.FAILED) {
        toastError(
          "Failed to process image, please try again",
          new Error("Image Processing Failed")
        );
        shouldClearState = true;
      } else if (
        progress.progress === Progress.PROCESSING ||
        progress.progress === Progress.READY
      ) {
        clearInterval(intervalRef.current);
        intervalRef.current = setInterval(() => {
          if (progressState.notifyReference) {
            // toastUpdate(progressState.notifyReference as string, {
            //   type: toast.TYPE.INFO,
            // });
            checkImageProgress();
          }
        }, 5000);
      }

      if (shouldClearState) {
        dispatchProgress({ type: "CLEAR" });
        intervalRef.current && clearInterval(intervalRef.current);
      }
    }
  };

  useEffect(() => {
    checkImageProgress();

    // return () => intervalRef.current && clearInterval(intervalRef.current);
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
