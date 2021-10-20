import { toast, ToastOptions } from "react-toastify";
import theme from "styles/colors";

const defaultToastOptions: ToastOptions = {
  position: "bottom-right",
  autoClose: 5000,
  hideProgressBar: true,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
  progress: undefined,
  style: {
    borderColor: theme.accent,
    borderWidth: 1,
    borderStyle: "solid",
  },
};

export const toastError = (msg: string, err: Error): void => {
  console.error(err);
  toast(msg, defaultToastOptions);
};
