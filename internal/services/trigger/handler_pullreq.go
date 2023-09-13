// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package trigger

import (
	"context"
	"fmt"

	"github.com/harness/gitness/events"
	"github.com/harness/gitness/internal/bootstrap"
	pullreqevents "github.com/harness/gitness/internal/events/pullreq"
	"github.com/harness/gitness/internal/pipeline/triggerer"
	"github.com/harness/gitness/types/enum"

	"github.com/drone/go-scm/scm"
)

func (s *Service) handleEventPullReqCreated(ctx context.Context,
	event *events.Event[*pullreqevents.CreatedPayload]) error {
	hook := &triggerer.Hook{
		Trigger:     enum.TriggerHook,
		Action:      enum.TriggerActionPullReqCreated,
		TriggeredBy: bootstrap.NewSystemServiceSession().Principal.ID,
		After:       event.Payload.SourceSHA,
	}
	err := s.augmentPullReqInfo(ctx, hook, event.Payload.PullReqID)
	if err != nil {
		return fmt.Errorf("could not augment pull request info: %w", err)
	}
	return s.trigger(ctx, event.Payload.SourceRepoID, enum.TriggerActionPullReqCreated, hook)
}

func (s *Service) handleEventPullReqReopened(ctx context.Context,
	event *events.Event[*pullreqevents.ReopenedPayload]) error {
	hook := &triggerer.Hook{
		Trigger:     enum.TriggerHook,
		Action:      enum.TriggerActionPullReqReopened,
		TriggeredBy: bootstrap.NewSystemServiceSession().Principal.ID,
		After:       event.Payload.SourceSHA,
	}
	err := s.augmentPullReqInfo(ctx, hook, event.Payload.PullReqID)
	if err != nil {
		return fmt.Errorf("could not augment pull request info: %w", err)
	}
	return s.trigger(ctx, event.Payload.SourceRepoID, enum.TriggerActionPullReqReopened, hook)
}

func (s *Service) handleEventPullReqBranchUpdated(ctx context.Context,
	event *events.Event[*pullreqevents.BranchUpdatedPayload]) error {
	hook := &triggerer.Hook{
		Trigger:     enum.TriggerHook,
		Action:      enum.TriggerActionPullReqBranchUpdated,
		TriggeredBy: bootstrap.NewSystemServiceSession().Principal.ID,
		After:       event.Payload.NewSHA,
	}
	err := s.augmentPullReqInfo(ctx, hook, event.Payload.PullReqID)
	if err != nil {
		return fmt.Errorf("could not augment pull request info: %w", err)
	}
	return s.trigger(ctx, event.Payload.SourceRepoID, enum.TriggerActionPullReqBranchUpdated, hook)
}

// augmentPullReqInfo adds in information into the hook pertaining to the pull request
// by querying the database.
func (s *Service) augmentPullReqInfo(
	ctx context.Context,
	hook *triggerer.Hook,
	pullReqID int64,
) error {
	pullreq, err := s.pullReqStore.Find(ctx, pullReqID)
	if err != nil {
		return fmt.Errorf("could not find pull request: %w", err)
	}
	hook.Title = pullreq.Title
	hook.Timestamp = pullreq.Created
	hook.AuthorLogin = pullreq.Author.UID
	hook.AuthorName = pullreq.Author.DisplayName
	hook.AuthorEmail = pullreq.Author.Email
	hook.Message = pullreq.Description
	hook.Before = pullreq.MergeBaseSHA
	hook.Target = pullreq.TargetBranch
	hook.Source = pullreq.SourceBranch
	// expand the branch to a git reference.
	hook.Ref = scm.ExpandRef(pullreq.SourceBranch, "refs/heads")
	return nil
}
