import { List, ListItem, ListItemIcon, Paper, Slide } from "@material-ui/core";
import _ from "lodash";
import * as React from "react";
import styled from "styled-components";
import Button from "./Button";
import Flex from "./Flex";
import Form from "./Form";
import FormCheckbox from "./FormCheckbox";
import Icon, { IconType } from "./Icon";
import Spacer from "./Spacer";
import Text from "./Text";

export type FilterConfig = { [key: string]: string[] };

const SlideContainer = styled.div`
  position: relative;
  width: 0px;
`;

const SlideWrapper = styled.div`
  position: absolute;
  right: 0;
  top: 0;
`;

/** Filter Bar Properties */
export interface Props {
  className?: string;
  /** the setState function for `activeFilters` */
  onFilterSelect: (val: FilterConfig) => void;
  /** Object containing column headers + corresponding filter options */
  filterList: FilterConfig;
  onClose?: () => void;
  open?: boolean;
}

function formStateToFilters(values) {
  const out = {};
  _.each(values, (v, k) => {
    const [key, val] = k.split(".");

    if (v) {
      const el = out[key];

      if (el) {
        el.push(val);
      } else {
        out[key] = [val];
      }
    }
  });

  return out;
}

function intialFormstate(cfg: FilterConfig) {
  return _.reduce(
    cfg,
    (r, vals, k) => {
      _.each(vals, (v) => {
        r[`${k}.${v}`] = true;
      });

      return r;
    },
    {}
  );
}

/** Form Filter Bar */
function UnstyledFilterDialog({
  className,
  onFilterSelect,
  filterList,
  onClose,
  open,
}: Props) {
  const onFormChange = ({ values }) => {
    if (onFilterSelect) {
      onFilterSelect(formStateToFilters(values));
    }
  };

  const initialState = intialFormstate(filterList);

  return (
    <SlideContainer>
      <SlideWrapper>
        <Slide direction="left" in={open} mountOnEnter unmountOnExit>
          <Paper elevation={4}>
            <Flex className={className + " filter-bar"} align start>
              <Spacer padding="medium">
                <Flex wide align between>
                  <Text size="extraLarge" color="neutral30">
                    Filters
                  </Text>
                  <Button variant="text" color="inherit" onClick={onClose}>
                    <Icon
                      type={IconType.ClearIcon}
                      size="large"
                      color="neutral30"
                    />
                  </Button>
                </Flex>
                <Form initialState={initialState} onChange={onFormChange}>
                  <List>
                    {_.map(filterList, (options: string[], header: string) => {
                      return (
                        <ListItem key={header}>
                          <Flex column>
                            <Text capitalize size="large" color="neutral30">
                              {header}
                            </Text>
                            <List>
                              {_.map(
                                options,
                                (option: string, index: number) => {
                                  return (
                                    <ListItem key={index}>
                                      <ListItemIcon>
                                        <FormCheckbox
                                          label=""
                                          name={`${header}.${option}`}
                                        />
                                      </ListItemIcon>
                                      <Text color="neutral30">
                                        {_.toString(option)}
                                      </Text>
                                    </ListItem>
                                  );
                                }
                              )}
                            </List>
                          </Flex>
                        </ListItem>
                      );
                    })}
                  </List>
                </Form>
              </Spacer>
            </Flex>
          </Paper>
        </Slide>
      </SlideWrapper>
    </SlideContainer>
  );
}

export default styled(UnstyledFilterDialog)`
  .MuiPopover-paper {
    min-width: 450px;
    border-left: 2px solid ${(props) => props.theme.colors.neutral30};
    padding-left: ${(props) => props.theme.spacing.medium};
  }
  .MuiListItem-gutters {
    padding-left: 0px;
  }
  .MuiCheckbox-root {
    padding: 0px;
  }
  .MuiCheckbox-colorSecondary {
    &.Mui-checked {
      color: ${(props) => props.theme.colors.primary};
    }
  }
`;