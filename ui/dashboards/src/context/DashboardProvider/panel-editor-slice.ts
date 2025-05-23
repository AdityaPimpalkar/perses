// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { Action, PanelEditorValues, PanelGroupId } from '@perses-dev/core';
import { StateCreator } from 'zustand';
import { getYForNewRow, getValidPanelKey } from '../../utils';
import { generateId, Middleware, createPanelDefinition } from './common';
import {
  PanelGroupSlice,
  PanelGroupItemId,
  PanelGroupDefinition,
  PanelGroupItemLayout,
  addPanelGroup,
  createEmptyPanelGroup,
} from './panel-group-slice';
import { PanelSlice } from './panel-slice';

/**
 * Slice that handles the visual editor state and actions for adding or editing Panels.
 */
export interface PanelEditorSlice {
  /**
   * Initial values for add panel if default panel kind is defined
   */
  initialValues?: Pick<PanelEditorValues, 'panelDefinition'>;

  /**
   * State for the panel editor when its open, otherwise undefined when it's closed.
   */
  panelEditor?: PanelEditorState;

  /**
   * Opens the editor for editing an existing panel by providing its layout coordinates.
   */
  openEditPanel: (panelGroupItemId: PanelGroupItemId) => void;

  /**
   * Opens the editor for adding a new Panel to a panel group.
   */
  openAddPanel: (panelGroupId?: PanelGroupId) => void;
}

export interface PanelEditorState {
  /**
   * Whether we're adding a new panel, or editing an existing panel.
   */
  mode: Action;

  /**
   * Initial values for the things that can be edited about a panel.
   */
  initialValues: PanelEditorValues;

  /**
   * Applies changes, but doesn't close the editor.
   */
  applyChanges: (next: PanelEditorValues) => void;

  /**
   * Close the editor.
   */
  close: () => void;
}

/**
 * Curried function for creating the PanelEditorSlice.
 */
export function createPanelEditorSlice(): StateCreator<
  // Actions in here need to modify both Panels and Panel Groups state
  PanelEditorSlice & PanelSlice & PanelGroupSlice,
  Middleware,
  [],
  PanelEditorSlice
> {
  // Return the state creator function for Zustand that uses the panels provided as intitial state
  return (set, get) => ({
    panelEditor: undefined,

    openEditPanel(panelGroupItemId): void {
      const { panels, panelGroups } = get();

      // Figure out the panel key at that location
      const { panelGroupId, panelGroupItemLayoutId: panelGroupLayoutId } = panelGroupItemId;
      const panelKey = panelGroups[panelGroupId]?.itemPanelKeys[panelGroupLayoutId];
      if (panelKey === undefined) {
        throw new Error(`Could not find Panel Group item ${panelGroupItemId}`);
      }

      // Find the panel to edit
      const panelToEdit = panels[panelKey];
      if (panelToEdit === undefined) {
        throw new Error(`Cannot find Panel with key '${panelKey}'`);
      }

      const editorState: PanelEditorState = {
        mode: 'update',
        initialValues: {
          groupId: panelGroupItemId.panelGroupId,
          panelDefinition: panelToEdit,
        },
        applyChanges: (next) => {
          set((state) => {
            state.panels[panelKey] = next.panelDefinition;

            // If the panel didn't change groups, nothing else to do
            if (next.groupId === panelGroupId) {
              return;
            }

            // Move panel to the new group
            const existingGroup = state.panelGroups[panelGroupId];
            if (existingGroup === undefined) {
              throw new Error(`Missing panel group ${panelGroupId}`);
            }

            const existingLayoutIdx = existingGroup.itemLayouts.findIndex((layout) => layout.i === panelGroupLayoutId);
            const existingLayout = existingGroup.itemLayouts[existingLayoutIdx];
            const existingPanelKey = existingGroup.itemPanelKeys[panelGroupLayoutId];
            if (existingLayoutIdx === -1 || existingLayout === undefined || existingPanelKey === undefined) {
              throw new Error(`Missing panel group item ${panelGroupLayoutId}`);
            }

            // Remove item from the old group
            existingGroup.itemLayouts.splice(existingLayoutIdx, 1);
            delete existingGroup.itemPanelKeys[panelGroupLayoutId];

            // Add item to the end of the new group
            const newGroup = state.panelGroups[next.groupId];
            if (newGroup === undefined) {
              throw new Error(`Could not find new group ${next.groupId}`);
            }

            newGroup.itemLayouts.push({
              i: existingLayout.i,
              x: 0,
              y: getYForNewRow(newGroup),
              w: existingLayout.w,
              h: existingLayout.h,
            });
            newGroup.itemPanelKeys[existingLayout.i] = existingPanelKey;
          });
        },
        close: () => {
          set((state) => {
            state.panelEditor = undefined;
          });
        },
      };

      // Open the editor with the new state
      set((state) => {
        state.panelEditor = editorState;
      });
    },

    openAddPanel(panelGroupId): void {
      // If a panel group isn't supplied, add to the first group or create a group if there aren't any
      let newGroup: PanelGroupDefinition | undefined = undefined;
      panelGroupId ??= get().panelGroupOrder[0];
      if (panelGroupId === undefined) {
        newGroup = createEmptyPanelGroup();
        newGroup.title = 'Panel Group';
        panelGroupId = newGroup.id;
      }

      const editorState: PanelEditorState = {
        mode: 'create',
        initialValues: {
          groupId: panelGroupId,
          panelDefinition: get().initialValues?.panelDefinition ?? createPanelDefinition(),
        },
        applyChanges: (next) => {
          const name = next.panelDefinition.spec.display.name;
          const panelKey = getValidPanelKey(name, get().panels);

          set((state) => {
            // Add a panel
            state.panels[panelKey] = next.panelDefinition;

            // Also add a panel group item referencing the panel
            const group = state.panelGroups[next.groupId];
            if (group === undefined) {
              throw new Error(`Missing panel group ${next.groupId}`);
            }
            const layout: PanelGroupItemLayout = {
              i: generateId().toString(),
              x: 0,
              y: getYForNewRow(group),
              w: 12,
              h: 6,
            };
            group.itemLayouts.push(layout);
            group.itemPanelKeys[layout.i] = panelKey;
          });
        },
        close: () => {
          set((state) => {
            state.panelEditor = undefined;
          });
        },
      };

      set((state) => {
        // Add the new panel group if one was created for the panel
        if (newGroup !== undefined) {
          addPanelGroup(state, newGroup);
        }

        // Open the editor with the new state
        state.panelEditor = editorState;
      });
    },
  });
}
