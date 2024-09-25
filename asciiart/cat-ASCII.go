package asciiart

// GetFrames returns the ASCII art frames of the cat animation (without whiskers)
func GetFrames() []string {
    return []string{
        `
  /\     /\  
 /  \___/  \ 
(  o     o  )
(     ^     )
 \_________/ 
  /       \  
 /         \ 
/_/|     |\_\
    `,
    `
  /\     /\  
 /  \___/  \ 
(  -     -  )
(     ^     )
 \_________/ 
  /       \  
 /         \ 
/_/|     |\_\
    `,
    `

  /\     /\  
 /  \___/  \ 
(  -     -  )
(     ^     )
 \_________/ 
 /         \ 
/_/|     |\_\
      `,
    }
}
