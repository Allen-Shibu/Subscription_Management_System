import { Show, SignIn, UserButton } from '@clerk/react'

export default function App() {
  return (
    <main>
      <Show when="signed-out">
        <SignIn />
      </Show>
      <Show when="signed-in">
        <UserButton />
      </Show>
    </main>
  )
}
