/* eslint-disable react/prop-types */
import React, { useState } from 'react'
import { useForm } from 'react-hook-form'
import { yupResolver } from '@hookform/resolvers/yup'
import * as Yup from 'yup'
import AddForm from './AddForm'

const AddFormContainer = () => {
    const [successful, setSuccessful] = useState(false)
    const [files, setFiles] = useState([])
    const [category, setCategory] = useState('')
    const [classification, setClassification] = useState('')

    const validationSchema = Yup.object().shape({
        classification: Yup.string().required('Classification is required'),
        category: Yup.string().required('Category is required'),
        superCategory: Yup.string().optional(),
        subCategory: Yup.string().optional(),
    })

    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(validationSchema),
    })

    const onSubmit = (data) => {
        if (Array.from(errors).length === 0) {
            const file = data.files
            setCategory(data.category)
            setClassification(data.classification)
            if (file) {
                setFiles(file)
            }
            setSuccessful(true)
        }
    }

    return (
        <div className="form-wrapper">
            <div>{`Classification: ${classification}\n`}</div>
            <div>{`Category: ${category}\n`}</div>
            <div>{`Files: ${files.length}\n`}</div>
            <AddForm
                successful={successful}
                onSubmit={onSubmit}
                errors={errors}
                register={register}
                handleSubmit={handleSubmit}
                reset={reset}
            />
        </div>
    )
}

export default AddFormContainer
