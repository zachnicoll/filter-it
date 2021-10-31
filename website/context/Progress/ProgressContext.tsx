import API from "api";
import { Progress } from "api/types";
import { toastError, toastSuccess, toastUpdate } from "common/toast";
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

  const intervalRef = useRef<any>(undefined);

  const checkImageProgress = async (): Promise<void> => {
    if (progressState.id) {
      let shouldClearState = false;

      const progress = await API.progress.get(progressState.id);

      if (progress === Progress.DONE) {
        toastSuccess("Image processed successfully!");
        shouldClearState = true;
      } else if (progress === Progress.FAILED) {
        toastError(
          "Failed to process image, please try again",
          new Error("Image Processing Failed")
        );
        shouldClearState = true;
      } else if (progress === Progress.PROCESSING) {
        intervalRef.current = setInterval(() => {
          if (progressState.notifyReference) {
            toastUpdate(progressState.notifyReference as string, {
              type: toast.TYPE.SUCCESS,
            });
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
