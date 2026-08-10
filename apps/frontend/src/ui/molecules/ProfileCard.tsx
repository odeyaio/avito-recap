import Avatar from "@mui/material/Avatar";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import { Link as RouterLink } from "react-router-dom";

import type { ProfileSummary } from "../../api/generated/model";

export interface ProfileCardProps {
  profile: ProfileSummary;
  to: string;
}

export function ProfileCard({ profile, to }: ProfileCardProps) {
  return (
    <Card variant="outlined">
      <CardActionArea component={RouterLink} to={to}>
        <CardContent>
          <Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
            <Avatar src={profile.avatarUrl} alt="">
              {profile.displayName.charAt(0)}
            </Avatar>
            <Stack spacing={0.5}>
              <Typography variant="subtitle1">
                {profile.displayName} · {profile.region}
              </Typography>
              {profile.teaser ? (
                <Typography variant="body2" color="text.secondary">
                  {profile.teaser}
                </Typography>
              ) : null}
            </Stack>
          </Stack>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
