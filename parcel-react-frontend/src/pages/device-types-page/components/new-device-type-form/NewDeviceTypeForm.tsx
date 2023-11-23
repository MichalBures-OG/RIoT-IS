import React, { useState, useCallback, useMemo, ChangeEvent } from 'react'
import { Button, FormControl, Grid, InputLabel, MenuItem, Select, SelectChangeEvent, TextField } from '@mui/material'
import styles from './NewDeviceTypeForm.module.scss'

interface NewDeviceTypeFormProps {
  createNewDeviceType: (denotation: string, parameters: { name: string; type: 'STRING' | 'NUMBER' | 'BOOLEAN' }[]) => Promise<void>
  anyLoadingOccurs: boolean
  anyErrorOccurred: boolean
}

const NewDeviceTypeForm: React.FC<NewDeviceTypeFormProps> = (props) => {
  const [denotation, setDenotation] = useState<string>('shelly1pro')
  const [parameters, setParameters] = useState<{ name: string; type: string }[]>([{ name: 'relay_0_temperature', type: 'NUMBER' }])

  const isFormDisabled: boolean = useMemo<boolean>(() => props.anyLoadingOccurs || props.anyErrorOccurred, [props.anyLoadingOccurs, props.anyErrorOccurred])

  const onSubmitHandler = useCallback(async () => {
    await props.createNewDeviceType(denotation, parameters as { name: string; type: 'STRING' | 'NUMBER' | 'BOOLEAN' }[])
  }, [denotation, parameters, props.createNewDeviceType])

  const onDenotationChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    setDenotation(e.target.value)
  }, [])

  const onParameterNameChange = (index: number) => (e: ChangeEvent<HTMLInputElement>) => {
    const newParameters = [...parameters]
    newParameters[index].name = e.target.value
    setParameters(newParameters)
  }

  const onParameterTypeChange = (index: number) => (e: SelectChangeEvent) => {
    const newParameters = [...parameters]
    newParameters[index].type = e.target.value
    setParameters(newParameters)
  }

  const addParameter = () => {
    setParameters([...parameters, { name: 'relay_0_temperature', type: 'NUMBER' }])
  }

  const deleteParameter = (index: number) => {
    const newParameters = [...parameters]
    newParameters.splice(index, 1)
    setParameters(newParameters)
  }

  return (
    <div className={styles.form}>
      <h2>Define a new device type</h2>
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={12}>
          <TextField fullWidth error={denotation.length === 0} label="Denotation" value={denotation} disabled={isFormDisabled} onChange={onDenotationChange} />
        </Grid>
        {parameters.map((parameter, index) => (
          <React.Fragment key={index}>
            <Grid item xs={4}>
              <TextField fullWidth error={parameter.name.length === 0} label={`Parameter ${index + 1} – Name`} value={parameter.name} disabled={isFormDisabled} onChange={onParameterNameChange(index)} />
            </Grid>
            <Grid item xs={4}>
              <FormControl fullWidth>
                <InputLabel id={`parameter-type-select-label-${index}`}>{`Parameter ${index + 1} – Type`}</InputLabel>
                <Select labelId={`parameter-type-select-label-${index}`} label={`Parameter ${index + 1} – Type`} value={parameter.type} onChange={onParameterTypeChange(index)} disabled={isFormDisabled}>
                  <MenuItem value={'STRING'}>STRING</MenuItem>
                  <MenuItem value={'NUMBER'}>NUMBER</MenuItem>
                  <MenuItem value={'BOOLEAN'}>BOOLEAN</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={4}>
              <Button disabled={isFormDisabled} onClick={() => deleteParameter(index)}>
                Delete this parameter
              </Button>
            </Grid>
          </React.Fragment>
        ))}
        <Grid item xs={12}>
          <Button disabled={isFormDisabled} onClick={addParameter}>
            Introduce next parameter
          </Button>
        </Grid>
        <Grid item xs={12}>
          <Button fullWidth disabled={isFormDisabled || denotation.length === 0 || parameters.some((p) => p.name.length === 0)} onClick={onSubmitHandler}>
            Submit
          </Button>
        </Grid>
      </Grid>
    </div>
  )
}

export default NewDeviceTypeForm
