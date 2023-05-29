import useAxios from 'axios-hooks'

export const useSubscribeNewsletter = (): [any, Function] => {
  const [
    { data, loading, error },
    execute
  ] = useAxios({url: "https://api.10x.pub/api-gateway/", method: "POST"}, {
     manual: true
  })

  const call = async (email: string): Promise<void> => {
    await execute({
      data: {
        operationName: 'SubscribeToNewsletter',
        variables: {email},
        query: 'mutation SubscribeToNewsletter($email: String!, $firstName: String, $lastName: String) {\n  subscribeToNewsletter(email: $email, firstName: $firstName, lastName: $lastName) {\n    ok\n    __typename\n  }\n}\n',
      }
    })
  }

  return [{ data, loading, error }, call];
}
