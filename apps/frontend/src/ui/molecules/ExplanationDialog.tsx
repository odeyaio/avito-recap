import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import Typography from "@mui/material/Typography";

import type { Evidence } from "../../api/generated/model";

export interface ExplanationDialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  explanation: string;
  evidence?: Evidence[];
}

export function ExplanationDialog({
  open,
  onClose,
  title,
  explanation,
  evidence,
}: ExplanationDialogProps) {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Typography gutterBottom>{explanation}</Typography>
        {evidence && evidence.length > 0 ? (
          <List dense disablePadding>
            {evidence.map((item) => (
              <ListItem key={item.metricCode} disableGutters>
                <ListItemText primary={item.description} />
              </ListItem>
            ))}
          </List>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
